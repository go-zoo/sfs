export SFSMASTERKEY="oskdghetyahsgdte"

rm sfs.dat
rm -rf test

mkdir test
echo "random file 1" > test/file1.txt
echo "random test file2" > test/file2.log
mkdir test/first
echo "test first file" > test/first/first1.txt
echo "another useless file" > test/first/useless.dat
mkdir test/first/dir
echo "deep file" > test/first/dir/deep.md
