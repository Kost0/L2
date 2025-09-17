#!/bin/bash
# Тестовый скрипт для нашего shell

echo "=== Test 1: Builtin commands ==="
pwd
echo "Hello world"
cd /tmp && pwd
cd ~ && pwd

echo "=== Test 2: External commands ==="
ls -la | wc
echo "test" | cat -n
ps aux | grep -v grep | head -2

echo "=== Test 3: Pipelines ==="
echo "line1\nline2\nline3" | wc -l
ps aux | grep $$ | wc -l

echo "=== Test 4: Redirections ==="
echo "content" > test_file.txt
cat test_file.txt
echo "more content" >> test_file.txt
cat < test_file.txt
wc -l < test_file.txt

echo "=== Test 5: Environment variables ==="
echo "Home: $HOME User: $USER"

echo "=== Test 6: Conditional execution ==="
true && echo "Success"
false || echo "Failed"
false && echo "This won't print"
true || echo "This won't print"

echo "=== Test 7: Error handling ==="
cd nonexistent_directory || echo "CD failed as expected"
unknown_command || echo "Command not found handled"

echo "=== Tests completed ==="