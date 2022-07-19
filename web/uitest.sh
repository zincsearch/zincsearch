#! /bin/sh

npm run test-once

rc=$?
if [ $rc -ne 0 ]; then
  echo "UI unit tests failed" >&2
  exit $rc
fi
