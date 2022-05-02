import json
import os
import sys

COMMIT_TYPES = set(["build", "ci", "docs", "feat", "fix", "perf", "refactor", "test", "style"])
MAX_COMMIT_MSG_LEN = 72

def commit_msg_len(msg):
    if len(msg) > MAX_COMMIT_MSG_LEN:
        print(f"PR title message '{msg}' too long. Max length is {MAX_COMMIT_MSG_LEN} characters.")
        sys.exit(1)

def message_format(msg):
    try:
        t = msg.split(":")[0].strip()
        m = msg.split(":")[1].strip()
    except:
        print(f"PR title format incorrect, should be 'type: message' ex: 'fix: a known issue'")
        sys.exit(1)

def commit_type(msg):
    t = msg.split(":")[0].strip()
    if t not in COMMIT_TYPES:
        print(f"PR type invalid. It needs to be one of {COMMIT_TYPES}. ")
        sys.exit(1)

def first_letter_case(msg):
    m = msg.split(":")[1].strip()
    if m[0].isupper():
        print(f"PR title message '{msg}' starts with upper case letter: '{m[0]}'")
        sys.exit(1)

def trailing_period(msg):
    if msg[-1] == ".":
        print(f"PR title message '{msg}' ends with a period. ")
        sys.exit(1)

def run_checks(msg):
    message_format(msg)
    commit_msg_len(msg)
    commit_type(msg)
    first_letter_case(msg)
    trailing_period(msg)

if __name__ == "__main__":

    with open(os.environ.get("GITHUB_EVENT_PATH")) as f:
        data = json.load(f)

        # is this a GitHub PR?
        if data["pull_request"]["title"]:
            msg = data["pull_request"]["title"]
            run_checks(msg)
        else:
            print("PR doesn't have a title, please add one that follows our guidelines.")
            sys.exit(1)
