import $ready from '../utils/ready';


$ready(async (Ajv, jsonata, { from, map, lastValueFrom }) => {
  console.log(Ajv);
  console.log(jsonata);
  console.log(from);
  console.log(map);
  console.log(lastValueFrom);

  const backstage = window.backstage;

  backstage.route('GET', '/ping', async (ctx, next) => {
    ctx.body = { code: 0, message: 'pong' };
    next();
  });
}, [ 'ajv', 'jsonata', 'rxjs' ]);
