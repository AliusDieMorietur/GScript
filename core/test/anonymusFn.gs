fn a(f, i) {
  f(i);
}

a(fn(i) {
  print i;
}, 3);

let b = fn(i) { print i; };

b(4);
