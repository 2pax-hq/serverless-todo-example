#!/bin/sh

set -e

# usage:
#
# #!/bin/sh
# dir=$(dirname "$0")
# fname=ListSessions
# source "$dir/../../invoke.sh"


srcdir=$dir/..

eventfile="$dir/example-event.json"

skipbuild=0

usage()
{
cat << EOF
usage: $0 options

Invoke function locally with the following options.

OPTIONS:
   -h      Show this message
   -e      Path to local event file to be used as input
   -s      Skip build step
EOF
}

while getopts "he:c:b:s" OPTION
do
  case $OPTION in
      h)
          usage
          exit 1
          ;;
      e)
          eventfile=$OPTARG
          ;;
      s)
          skipbuild=1
          ;;
      ?)
          usage
          exit
          ;;
  esac
done

if [[ $skipbuild == 0 ]]
then
  echo "compiling..."
  "$dir/build.sh"
fi


sam local invoke "$fname" \
  -t "$srcdir/../../cloudformation.yaml" \
  -e "$eventfile"
