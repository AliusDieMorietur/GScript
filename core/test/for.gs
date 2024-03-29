for (let a = 5; a != 0; a = a - 1) {
  print a;
}

for (let a = 10; a != 0; a = a - 1) {
  if (a == 8) {
    continue;
  }
  if (a < 5) {
    break;
  }
  print a;
}
