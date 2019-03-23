# Name
treenum - A highly extensible, directory-based host enumeration tool.
# Synopsis
treenum *hostname* *script-folder* [--nodetach]
# Description
Treenum is an extensible port-scanning tool that outputs a neat hierarchy of
"branch" directories for a target host detailing open TCP/UDP ports.  Treenum
can be thought of as a programmable extension to nmap, which encodes its output
into a set of easily parsed and understood directories, in contrast to the
meriod of tools which push everything to stdout or their own xml/json/"custom"
file formats.

Treenum was made to bridge the gap between effective note-taking software like
CherryTree and port scanning software like nmap, and to allow users to encode
their host enumeration strategy into set of scripts and tools.
