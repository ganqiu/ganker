# init
mkdir lower1
mkdir lower2
mkdir upper
mkdir work
mkdir merged
echo I\'m file1, belong to lower1 > lower1/file1
echo I\'m file2, belong to lower2 > lower2/file2
echo I\'m file3, belong to upper > upper/file3
sudo mount -t overlay overlay2 -olowerdir=lower1:lower2,upperdir=upper,workdir=work merged/

#try to change file1
echo I\'m file1, belong to lower1, but I\'m changed > merged/file1
cat merged/file1
cat lower1/file1
tree 

# try to delete file2
rm merged/file2.txt
tree
cat merged/*
cat lower2/file2.txt
ll upper/file2.txt