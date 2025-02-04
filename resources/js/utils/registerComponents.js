import { createApp } from 'vue';

export function autoRegisterComponents(components) {
  Object.entries(components).forEach(([name, component]) => {
    document.querySelectorAll(name).forEach((el) => {
      const app = createApp(component);

      const props = Object.fromEntries(
        Array.from(el.attributes)
          .filter(attr => attr.name.startsWith('data-'))
          .map(attr => [attr.name.replace('data-', ''), attr.value])
      );

      app.mount(el, props);
    });
  });
}