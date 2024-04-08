fn a(f, i) {
  let a = 3;
  print a;
  f(i);
}

fn b(i) {
  print i;
}

a(b, 2);
