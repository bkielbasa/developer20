b64=`cat baseImg.png | base64`

for dir in content/post/*
do
    if test -f "$dir/featured.jpg"; then
        continue
    fi

    if test -f "$dir/featured.png"; then
        continue
    fi

    echo "$dir doesn't have the featured img. Creating..."

    title=`grep 'title' $dir/index.md | cut -c 8-`
    if [[ ${title:0:1} == "\"" ]] ; then 
        title=$(sed 's/.\{1\}$//' <<< "$title")
        title=$(sed 's/^.\{1\}//' <<< "$title")
    fi

    echo $title

    curl -X POST -H "Content-Type: application/json" \
    -d "{\"title\": \"$title\", \"baseImg\": \"$b64\"}" \
    https://imgs.developer20.com/blog-post --output $dir/featured.png
done
