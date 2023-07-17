
deploy:
	hugo
	aws s3 cp public s3://developer20.com/ --recursive --acl public-read

