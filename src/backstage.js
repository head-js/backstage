import * as logger from './utils/logger';
import requirejs from './utils/requirejs';


logger.debug('backstage.js');


requirejs('js/frontstage.js');
