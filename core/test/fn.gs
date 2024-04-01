fn count(n) {
  if (n > 1) count(n - 1);
  print n;
}

print count; 
count(3);

fn helloWorld(end) {
  print "Hello, world! " + end;
}

helloWorld("42");

fn sum(a, b, c) {
  return a + b + c;
}

print sum(1,2,3);

fn fib(n) {
  if (n <= 1) return n;
  return fib(n - 2) + fib(n - 1);
}

for (let i = 0; i < 20; i = i + 1) {
  print fib(i);
}
