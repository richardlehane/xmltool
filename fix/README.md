Package xmltool/fix repairs invalid XML files.

Written to clean up generic XML files generated by databases & that contain invalid entities.

Example:

    reader := strings.NewReader("<dodgy><hello>Richa&rd</hello><richard.lehane@gmail.com></dodgy>")
    writer := new(bytes.Buffer)
    error := fixxml.Fixxml(reader, writer)
    if writer.String() == "<dodgy><hello>Richa&amp;rd</hello>&lt;richard.lehane@gmail.com&gt;</dodgy>" {
        fmt.Print("Thanks, I've been Fixxmled!")
    }