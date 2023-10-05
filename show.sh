raw=$(git show -s --date=unix --format="%cd")
hex=$(git show -s --format="%H")
python3 t.py $raw $hex