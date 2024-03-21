import argparse

from typing import List
from os import path, listdir
import subprocess
import apiTesting
import utils.globalConfig as conf


parser = argparse.ArgumentParser(
    description='Run a regression test on InterinGo project.')
parser.add_argument('--debug', dest='debug',
                    action='store_true',
                    help='Debug mode, output will be more verbose and the sript suppose to be run in `script/` directory')
parser.add_argument('--input', dest='inDir',
                    action='store',
                    default="test/",
                    help='Path that contain all the input file, default "test/"')
parser.add_argument('--output', dest='outDir',
                    action='store',
                    default="test/result/",
                    help='Path that contain all outputfile, default "test/result/"')
parser.add_argument('--outEncode', dest='outEncode',
                    action='store',
                    default="ascii",
                    help='Set the default encode for command output')
args = parser.parse_args()

conf.DEBUG = args.debug
conf.IN_DIR = args.inDir
conf.OUT_DIR = args.outDir
conf.OUT_ENCODE = args.outEncode

# The script intent to be run in project root directory rather than in script/
# Change this flag to True only when you need to debug it (inside script/)
if conf.DEBUG:
    from os import chdir
    chdir("..")


def validateArgs():
    if not path.isdir(args.inDir):
        raise ValueError("Input directory is not a valid path")
    if not path.isdir(args.outDir):
        raise ValueError("Output directory is not a valid path")


def buildCommand(*args, **kwargs) -> str:
    """Helper to build test command

    Args:
        input: the input file path
        *args: additional single flag without `-`
        **kwargs: additional value flag without `-`

    Returns:
        A string that contain build command with all the flag
    """

    base = f"{conf.EXEC_PATH}"
    singleFlag = " ".join([f"-{i}" for i in args])
    valueFlag = " ".join(
        [f"-{flag}={value}" for (flag, value) in kwargs.items()])
    return " ".join([base, singleFlag, valueFlag])


def buildFileModeCommand(inputPath: str, *args, **kwargs):
    """Helper to build File mode command

    Args:
        inputPath: `*.iig` file path
        *args: Others single flag
        **kwargs: Other value flag

    Returns:
        A string that contain build command with key flag
    """
    return buildCommand(f=inputPath, *args, **kwargs)


def getInputFileList() -> List[str]:
    """ Return list of input file """

    # Dirty trick as I know there can't be any other file in there with `.iig`
    # in its name
    return [i for i in listdir(conf.IN_DIR) if ".iig" in i]


def checkOutFile(outputFilePath: str, result: bytes):
    """recheck output file with a command result to see if it match

    Args:
        outputFilePath: Output file path
        result: Command new output

    Raises:
        OSError: Raise if file can't not access
    """

    try:
        fout = open(outputFilePath, 'rb')
        oldOutput = fout.read()
        fout.close()
    except OSError as e:
        raise OSError(
            f"Can't open output file, please check environment. Skipping {outputFilePath}", e)

    diffcheck = False
    for i, b in enumerate(oldOutput):
        if result[i] != b:
            diffcheck = True
            break

    if diffcheck:
        print(f"Output change, please recheck {outputFilePath} manually")


def evaluateIigFile(inputFileName: str):
    command = buildFileModeCommand(inputFileName)
    out = subprocess.Popen(command, shell=True, stdout=subprocess.PIPE).stdout
    if out is not None:
        return out.read()
    return b""


def fileModeTest():
    if conf.DEBUG:
        print("INFO: Start regression check File mode result")
    testFiles = getInputFileList()
    for fn in testFiles:
        if conf.DEBUG:
            print(f"INFO: Check {fn} test file output ...")

        fullPathInput = path.join(conf.IN_DIR, fn)
        fullPathOutput = path.join(conf.OUT_DIR, fn)[:-4] + '.out'

        result = evaluateIigFile(fullPathInput)
        checkOutFile(fullPathOutput, result)
    print(f"PASS: Checked {len(testFiles)} input file")


def goTest():
    if conf.DEBUG:
        print("Start go native test...")
    out = subprocess.Popen("go test ./...",
                           shell=True, stdout=subprocess.PIPE).stdout

    if out is None:
        return

    goTestOutput = out.read().decode(conf.OUT_ENCODE)
    lines = goTestOutput.split("\n")
    for line in lines:
        if line == "":
            continue
        if line[0] == '?':
            print(
                f"WARNING: Module {line.split()[1]} Don't have native test, consider adding it")

    # Temprorary implement until parse is finish
    if lines is not None and lines[-2] == "FAIL":
        print("FAIL: go native test fail")
    else:
        print("PASS: go native test done")


if __name__ == "__main__":
    validateArgs()
    fileModeTest()
    apiTesting.serverTest(conf.EXEC_PATH, conf.IN_DIR, conf.OUT_DIR)
    goTest()
