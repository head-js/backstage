export function ignore(message, extra1 = '') {
  console.log('%c%s', 'color: #727272', message, extra1);
}


export function debug(message, extra1 = '', extra2 = '', extra3 = '') {
  console.log('%c%s', 'background-color: #f0f9ff', message, extra1, extra2, extra3);
}
