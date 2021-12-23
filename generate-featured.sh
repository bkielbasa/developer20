b64=`cat baseImg.png | base64`

echo "Creating..."
for dir in content/post/*
do
    if test -f "$dir/featured.jpg"; then
        continue
    fi

    title=`grep 'title' $dir/index.md | head -n 1 | cut -c 8-`
    if [[ ${title:0:1} == "\"" ]] ; then 
        title=$(sed 's/.\{1\}$//' <<< "$title")
        title=$(sed 's/^.\{1\}//' <<< "$title")
    fi

    curl -X POST -H "Content-Type: application/json" \
    -d "{\"title\": \"$title\", \"baseImg\": \"$b64\"}" \
    https://imgs.developer20.com/blog-post --output $dir/featured.png &> /dev/null
done
