rm submissions.csv
rm submissions.csv.gz
wget https://s3-ap-northeast-1.amazonaws.com/kenkoooo/submissions.csv.gz
date > date.txt
gunzip submissions.csv.gz 
