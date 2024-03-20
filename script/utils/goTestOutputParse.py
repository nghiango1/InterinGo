from typing import Dict
from enum import Enum
import json
from datetime import datetime


class GoTestLine:
    """Helper class for Go test result json parse"""

    def __init__(self, value: Dict):
        """Create new gotest object fomr

        Args:
            value: Output line from go test command with `-json` flag have been parse to Dict
        """

        self.time: str = datetime.fromisoformat(value["Time"])


class ActionType(Enum):
    OUTPUT = "output"
    PASS = "pass"
    RUN = "run"
    FAIL = "false"
    START = "start"


class Action:
    """Hacky way to parse the Action"""

    def __init(self, value: str, **kwargs):
        try:
            self.type = ActionType[value]
        except KeyError as e:
            print(e)
            raise KeyError(f"Action string type isn't supported, got {value}")

        # {"Time":"2024-03-20T20:53:39.371651163+07:00","Action":"start","Package":"interingo"}
        self.package = kwargs["Package"]

        if self.type == ActionType.START:
            return

        self.test = kwargs["Test"]
        if self.type == ActionType.PASS or self.type == ActionType.FAIL:
            self.elapsed = kwargs["Elapsed"]
        elif self.type == ActionType.RUN:
            return
        elif self.type == ActionType.OUTPUT:
            self.output = kwargs["Output"]


class PackageTestResult:
    def __init__(self, package: str):
        self.package = package
        self.test: Dict[str, TestResult] = {}

    def addTest(self, test: str):
        self.test[1]


class TestResult:
    def __init__(self, package: str, test: str):
        self.package = package
        self.test = test
        self.status = None
        self.actions = []
        self.startTime = None
        self.elapsedTime = None

    def addAction(self, action: Action):
        self.actions.append(action)


class GoTestOut:
    """Helper class for Go test result json parse

    Attributes:
        time: time
    """

    def __init__(self, lines: bytes):
        """Create new gotest object fomr

        Args:
            line: Output line from go test command with `-json` flag

        Raises:
            ValueError: Parsing input error
        """
        self.test = []
        self.time = []
        self.action = []
        self.pakage = []
        self.test = []
        self.output = []

        for line in str(lines).split:
            parseLine = None
            try:
                parseDict: Dict = json.load(line)
                parseLine
            except ValueError as e:
                print(e)
                raise ValueError(
                    "Can't parse go test output, please recheck test command")
