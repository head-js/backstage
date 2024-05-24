export default function requirejs(src) {
  const s = document.createElement('script');
  s.src = chrome.runtime.getURL(src);
  // s.type = 'module';
  s.onload = () => s.remove();
  (document.head || document.documentElement).append(s);
}
