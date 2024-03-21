import argparse

from typing import List, Tuple
from os import path, listdir
import subprocess


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

DEBUG: bool = args.debug
TEST_DIR: str = args.inDir
OUT_DIR: str = args.outDir
OUTPUT_ENCODE: str = args.outEncode
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


def buildCommand(*args: str, **kwargs) -> str:
    """Helper to build test command

    Args:
        input: the input file path
        *args: additional single flag string without `-`
        **kwargs: additional value flag without `-`

    Returns:
        A string that contain build command with all the flag
    """

    base = f"{TEST_CMD}"
    singleFlag = " ".join([f"-{i}" for i in args])
    valueFlag = " ".join(
        [f"-{flag}={value}" for (flag, value) in kwargs.items()])
    return " ".join([base, singleFlag, valueFlag])


def buildDevServerCommand(listenAddress: str = "127.0.0.1:0", *args: str, **kwargs):
    return buildCommand("s", l=listenAddress, *args, **kwargs)


def getInputFileList() -> List[str]:
    """ Return list of input file """

    # Dirty trick as I know there can't be any other file in there with `.iig`
    # in its name
    return [i for i in listdir(TEST_DIR) if ".iig" in i]


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
            f"ERROR: Can't open output file, please check environment. Skipping {outputFilePath}", e)

    diffcheck = False
    for i, b in enumerate(oldOutput):
        if result[i] != b:
            diffcheck = True
            break

    if diffcheck:
        print(f"FAIL: Output change, please recheck {outputFilePath} manually")
    else:
        if DEBUG:
            print("DEBUG: Server response match provided output")


class REPLServer:
    def __init__(self):
        self.serverURL = REPLServer.getFreeLocalAdr()
        self.process = None
        self.startREPLServer(self.serverURL)

    def startREPLServer(self, listenAddress: str):
        """This not epecte to be fail, but our reserved Listen port can be taken in the span time of python running these code

        Args:
            listenAddress: Address that REPL server listen
        """
        if DEBUG:
            print("DEBUG: Start server subprocess")
        if self.process is not None:
            if DEBUG:
                print(
                    f"DEBUG: Server subprocess already start, pid={self.process.pid}, URL={self.serverURL}")
            return
        command = buildDevServerCommand(listenAddress)
        self.process = subprocess.Popen(
            command, shell=True, stdout=subprocess.PIPE)
        self.serverURL = listenAddress

    @staticmethod
    def getFreeLocalAdr():
        """Get a free local TCP port for REPL server

        Returns: A local address `127.0.0.1:{port}` that is currently free
        """
        import socket

        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.bind(('', 0))  # Bind to a free port provided by the host.
        port = s.getsockname()[1]  # Retrieve the port number.
        s.close()  # Close the socket, making the port free again.

        return f"127.0.0.1:{port}"

    replEndpoint = "/api/evaluate"

    def sendInputToREPLEndpoint(self, input) -> Tuple[bytes, int]:
        import requests

        body = {'repl-input': input}
        r = requests.post(
            url=f"http://{self.serverURL}{REPLServer.replEndpoint}", data=body)
        if r.status_code == 200:
            return r.content, r.status_code
        else:
            print(
                f"FAIL: Server not responding, got status code: {r.status_code}")
            return b"", r.status_code

    def closeREPLServer(self):
        if self.process is not None:
            # Our server isn't need clean up so 9 is better to ensure the process end
            # self.process.send_signal(9)
            self.process.send_signal(15)


def getInputFromFile(inputFilePath: str):
    try:
        fout = open(inputFilePath, 'rb')
        oldOutput = fout.read()
        fout.close()
        return oldOutput
    except OSError as e:
        raise OSError(
            f"Can't open input file, please check environment. Skipping {inputFilePath}", e)


def serverTest():
    """Check the REPL API result spawn by ./interingo -s"""

    if DEBUG:
        print("DEBUG: Start API test on REPL server")
    testFiles = getInputFileList()
    replServer = REPLServer()

    for fn in testFiles:
        if DEBUG:
            print(f"DEBUG: Check {fn} test file - sending to API...")

        fullPathInput = path.join(TEST_DIR, fn)
        fullPathOutput = path.join(OUT_DIR, fn)[:-4] + '.out'

        input = getInputFromFile(fullPathInput)
        result, statusCode = replServer.sendInputToREPLEndpoint(input)
        if DEBUG:
            print(f"DEBUG: Server response {result[:20]}...")
        if statusCode == 200:
            checkOutFile(fullPathOutput, result)

    replServer.closeREPLServer()
    print(f"PASS: Checked {len(testFiles)} input file")


if __name__ == "__main__":
    validateArgs()
    serverTest()
