import argparse

from typing import List
from os import path, listdir
import subprocess


parser = argparse.ArgumentParser(
    description='Generating test result from *iig file.')
parser.add_argument('--debug', dest='debug',
                    action='store_true',
                    help='Debug mode, output will be more verbose and the sript suppose to be run in `script/` directory')
parser.add_argument('--input', dest='inDir',
                    action='store',
                    default="test/",
                    help=f'Path that contain all the input file, default "test/"')
parser.add_argument('--output', dest='outDir',
                    action='store',
                    default="test/result/",
                    help='Path that will contain all outputfile, default "test/result/"')
args = parser.parse_args()

DEBUG = args.debug
TEST_DIR = args.inDir
OUT_DIR = args.outDir
TEST_CMD = "./interingo"
# The script intent to be run in project root directory rather than in script/
# Change this flag to True only when you need to debug it (inside script/)
if DEBUG:
    from os import chdir
    chdir("..")


def validateArgs():
    if not path.isdir(args.inDir):
        raise ValueError("Input directory is not a valid path")
    if not path.isdir(args.outDir):
        raise ValueError("Output directory is not a valid path")


def buildCommand(input: str, *args, **kwargs) -> str:
    """Helper to build test command

    Args:
        input: the input file path
        *args: additional single flag without `-`
        **kwargs: additional value flag without `-`

    Returns:
        A string that contain build command with all the flag
    """

    base = f"{TEST_CMD} -f {input}"
    singleFlag = " -".join([i for i in args])
    valueFlag = " -".join(
        [f"{flag}={value}" for (flag, value) in kwargs.items()])
    return " ".join([base, singleFlag, valueFlag])


def getInputFileList() -> List[str]:
    """ Return list of input file """

    # Dirty trick as I know there can't be any other file in there with `.iig`
    # in its name
    return [i for i in listdir(TEST_DIR) if ".iig" in i]


def recheckOutFile(outputFilePath: str, result: str):
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


def writeOutFile(outputFilePath: str, result: str):
    """white command result into output file

    Args:
        outputFilePath: output file path
        result: command result
    """
    # If theres output file already exist, just recheck output file instead
    isOldOutputExist = path.isfile(outputFilePath)

    if isOldOutputExist:
        recheckOutFile(outputFilePath, result)
        return

    # Else we white generated output
    fout = open(outputFilePath, 'wb')
    fout.write(result)
    fout.close()


def evaluateIigFile(inputFileName: str):
    command = buildCommand(fullPathInput)
    return subprocess.Popen(command, shell=True,
                            stdout=subprocess.PIPE).stdout.read()


if __name__ == "__main__":
    validateArgs()
    testFiles = getInputFileList()
    for fn in testFiles:
        if DEBUG:
            print(f"Genereting {fn} test file out put ...")

        fullPathInput = path.join(TEST_DIR, fn)
        fullPathOutput = path.join(OUT_DIR, fn)[:-4] + '.out'

        result = evaluateIigFile(fullPathInput)
        writeOutFile(fullPathOutput, result)
    print(f"Done generate form {len(testFiles)} input file")
