#!/bin/sh

if [ $# -ne 3 ]; then
  echo "Usage $0 <worklog_url> <month> <dest-dir>"
  exit 2
fi

worklog_url=$1
month=$2
destdir=$3
lowermonth="$(echo "$month" | tr '[:upper:]' '[:lower:]')"

workdir=$(mktemp -d)
export WORKLOG="$workdir/worklog.txt"

wget "$worklog_url" -O "$WORKLOG"

worklog filter "$month" | worklog invoice >"$workdir/invoice.html"
worklog filter "$month" | worklog fmt html >"$workdir/worklog.html"

chromium --headless --print-to-pdf-no-header --disable-gpu --print-to-pdf="$destdir/invoice_${lowermonth}.pdf" "$workdir/invoice.html"
chromium --headless --print-to-pdf-no-header --disable-gpu --print-to-pdf="$destdir/worklog_${lowermonth}.pdf" "$workdir/worklog.html"

rm -r "$workdir"
