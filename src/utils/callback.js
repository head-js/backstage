function seq() {
  return Date.now() + '-' + Math.random().toString(36).substring(2);
}


export default function Callback() {
  const promised = {};

  return function callback(fn, resolved, rejected) {
    if (typeof fn === 'function') {
      const id = seq();
      return new Promise((resolve, reject) => {
        promised[id] = { resolve, reject };
        fn(id);
      });
    } else { // typeof fn === 'string'
      const { resolve, reject } = promised[fn];
      if (rejected) {
        reject(rejected);
      } else {
        resolve(resolved);
      }
      delete promised[fn];
    }
  };
}
