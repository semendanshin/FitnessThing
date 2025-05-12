import fs from 'fs';
import path from 'path';
import Handlebars from 'handlebars';

const templateCache = new Map<string, HandlebarsTemplateDelegate>();

function loadTemplate(name: string): HandlebarsTemplateDelegate {
  if (templateCache.has(name)) {
    return templateCache.get(name)!;
  }
  const filePath = path.join(__dirname, '..', 'templates', `${name}.html`);
  const source = fs.readFileSync(filePath, 'utf8');
  const template = Handlebars.compile(source);
  templateCache.set(name, template);
  return template;
}

export function renderTemplate(name: string, context: Record<string, any>): string {
  const template = loadTemplate(name);
  return template(context);
}
