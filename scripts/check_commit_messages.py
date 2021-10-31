import json
import os
import sys

COMMIT_TYPES = set(["build", "ci", "docs", "feat", "fix", "perf", "refactor", "test"])
MAX_COMMIT_MSG_LEN = 72


def commit_msg_len(msg):
    if len(msg) > MAX_COMMIT_MSG_LEN:
        print(f"Commit message '{msg}' too long.")
        sys.exit(1)


def commit_type(msg):
    t = msg.split(":")[0].strip()
    if t not in COMMIT_TYPES:
        print(f"Commit type invalid. It needs to be one of {COMMIT_TYPES}")
        sys.exit(1)


def first_letter_case(msg):
    m = msg.split(":")[1].strip()
    if m[0].isupper():
        print(f"Commit message '{msg}' starts with upper case letter.")
        sys.exit(1)


def trailing_period(msg):
    if msg[-1] == ".":
        print(f"Commit message '{msg}' ends with a period.")
        sys.exit(1)


if __name__ == "__main__":

    with open(os.environ.get("GITHUB_EVENT_PATH")) as f:
        data = json.load(f)
        for c in data["commits"]:
            msg = c["message"]
            print(msg)
            commit_msg_len(msg)
            commit_type(msg)
            first_letter_case(msg)
            trailing_period(msg)
