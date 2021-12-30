b64=`cat baseImg.png | base64`

echo "Creating for blog posts..."
for dir in content/post/*
do
    title=`grep 'title' $dir/index.md | head -n 1 | cut -c 8-`
    if [[ ${title:0:1} == "\"" ]] ; then 
        title=$(sed 's/.\{1\}$//' <<< "$title")
        title=$(sed 's/^.\{1\}//' <<< "$title")
    fi

    curl -X POST -H "Content-Type: application/json" \
    -d "{\"title\": \"$title\", \"baseImg\": \"$b64\"}" \
    https://imgs.developer20.com/blog-post --output $dir/featured.png &> /dev/null
done

echo "Creating for reviews..."
for dir in content/reviews/*
do
    title=`grep 'title' $dir/index.md | head -n 1 | cut -c 8-`
    if [[ ${title:0:1} == "\"" ]] ; then 
        title=$(sed 's/.\{1\}$//' <<< "$title")
        title=$(sed 's/^.\{1\}//' <<< "$title")
    fi

    curl -X POST -H "Content-Type: application/json" \
    -d "{\"title\": \"$title\", \"baseImg\": \"$b64\"}" \
    https://imgs.developer20.com/blog-post --output $dir/featured.png &> /dev/null
done