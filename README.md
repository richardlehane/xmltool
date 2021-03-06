Cleans up generic XML files generated by databases. Also audits xml files, counting occurrence of elements.
If given a directory, will do a recursive walk, fixing or auditing any files with a ".xml" extension.

Examples:

    ./xmltool -fix bad.xml > good.xml
    ./xmltool -audit good.xml
    ./xmltool -fix DIR_CONTAINING_BAD_XML_FILES -outdir ~/Good
    ./xmltool -audit ~/Good -html > report.html

Compile with `go install` or grab a binary from the Github releases.

[![Build Status](https://travis-ci.org/richardlehane/xmltool.png?branch=master)](https://travis-ci.org/richardlehane/xmltool)