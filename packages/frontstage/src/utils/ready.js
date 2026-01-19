const SOURCE = '@head/backstage';


async function importFromBackstageVendors(vendors) {
  const loaded = [];
  for (let i = 0; i < vendors.length; i += 1) {
    const v = vendors[i];
    switch (v) {
      case 'ajv':
        const { default: ajv } = (await __backstagevendors__.get('./ajv'))();
        loaded.push(ajv);
        break;
      case 'jsonata':
        const { default: jsonata } = (await __backstagevendors__.get('./jsonata'))();
        loaded.push(jsonata);
        break;
      case 'rxjs':
        const { default: rxjs } = (await __backstagevendors__.get('./rxjs'))();
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
  if (window.backstage && window.backstage.route) {
    const loaded = await importFromBackstageVendors(vendors);
    fn(...loaded);
  } else {
    window.addEventListener('message', async ({ /* type, source, origin */ data }) => {
      if (data.source !== SOURCE) {
        return false;
      }

      if (data.type === 'READY') {
        const loaded = await importFromBackstageVendors(vendors);
        fn(...loaded);
      }
    });
  }
}
