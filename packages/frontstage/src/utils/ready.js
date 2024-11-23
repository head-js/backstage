const SOURCE = '@head/backstage';


async function importFromBackstageVendors(vendors) {
  const loaded = [];
  for (let i = 0; i < vendors.length; i += 1) {
    const v = vendors[i];
    switch (v) {
      case 'ajv':
        const { default: ajv } = await import('backstage/ajv');
        loaded.push(ajv);
        break;
      case 'jsonata':
        const { default: jsonata } = await import('backstage/jsonata');
        loaded.push(jsonata);
        break;
      case 'rxjs':
        const { default: rxjs } = await import('backstage/rxjs');
        loaded.push(rxjs);
        break;
      default:
        loaded.push(null);
        break;
    }
  }
  // console.log(loaded);
  return loaded;
}


export default async function $ready(fn, vendors) {
  const loaded = await importFromBackstageVendors(vendors);

  if (window.backstage && window.backstage.route) {
    fn(...loaded);
  } else {
    window.addEventListener('message', ({ /* type, source, origin */ data }) => {
      if (data.source !== SOURCE) {
        return false;
      }

      if (data.type === 'READY') {
        fn(...loaded);
      }
    });
  }
}
