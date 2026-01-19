export default function requirejs(src) {
  return new Promise((resolve, reject) => {
    const s = document.createElement('script');
    s.src = chrome.runtime.getURL(src);
    // s.type = 'module';
    s.onload = () => {
      s.remove();
      resolve();
    };
    s.onerror = () => {
      s.remove();
      reject(new Error(`Failed to load script: ${src}`));
    };
    (document.head || document.documentElement).append(s);
  });
}
