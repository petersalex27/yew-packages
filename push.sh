# does the following sequence of commands:
# in : git add *
# in : git commit -m "$MSG"
# in : git push $REPO $REF
# in : sh show.sh
# out: <command_output>
# out: @<pseudo_version_for_go_get>

main_script() {
  if [ -z "$1" ]; then
    echo "usage: sh push.sh <commit_message> <repo> [ref]"
    echo "                  ^^^^^^^^^^^^^^^^ missing commit message"
    exit
  fi

  MSG="$1"

  if [ -z $2 ]; then
      echo "usage: sh push.sh <commit_message> <repo> [ref]"
      echo "                                   ^^^^^^ missing remote repository"
      exit
  fi

  REPO="$2"
  REF=""

  if [ -z $3 ]; then
    REF=$(git rev-parse --abbrev-ref HEAD)
    echo "using current ref=$REF"
  else
    REF=$3
  fi

  git add --all && git commit -m "$MSG" && git push "$REPO" "$REF" && sh show.sh
}

main_script "$@"
exit