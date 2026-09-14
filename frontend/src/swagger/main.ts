import SwaggerUIBundle from 'swagger-ui-dist/swagger-ui-bundle.js';
import 'swagger-ui-dist/swagger-ui.css';
import { apiEndpoint } from '../services/endpoint';
import { searchPlugin } from './search';

SwaggerUIBundle({
  url: `${apiEndpoint}/swagger/doc.json`,
  dom_id: '#swagger-ui',
  deepLinking: true,
  filter: true,
  plugins: [searchPlugin],
});
