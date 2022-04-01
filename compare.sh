#!/bin/bash
echo "Test #1"
diff -y <(go run . < samples/input1.txt) samples/output1.txt

echo "Test #2"
diff -y <(go run . < samples/input2.txt) samples/output2.txt

echo "Test #3"
diff -y <(go run . < samples/input3.txt) samples/output3.txt

echo "Test #4"
diff -y <(go run . < samples/input4.txt) samples/output4.txt

echo "Test #5"
diff -y <(go run . < samples/input5.txt) samples/output5.txt