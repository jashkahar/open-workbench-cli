<template>
  <div class="container">
    <h1>Welcome to {{.ProjectName}}</h1>
    <p>
      Owned by: <strong>{{.Owner}}</strong>
    </p>

    <div class="features">
      <h2>Features</h2>
      <ul>
        <li>⚡️ Nuxt 3 - The Vue.js framework for production</li>
        <li>🔷 TypeScript - Type-safe development</li>
        <li>🎨 Tailwind CSS - Utility-first CSS framework</li>
        <li>🧪 Vitest - Fast unit testing</li>
        <li>📦 ESLint - Code linting and formatting</li>
        <li>🐳 Docker - Containerization support</li>
        <li>📱 PWA - Progressive Web App features</li>
        <li>🔧 Auto-imports - Automatic component imports</li>
      </ul>
    </div>

    <div class="getting-started">
      <h2>Getting Started</h2>
      <p>Edit <code>pages/index.vue</code> to customize this page.</p>
      <p>
        Check out the
        <a href="https://nuxt.com/docs" target="_blank">Nuxt documentation</a>
        to learn more.
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
// This page uses the Composition API with TypeScript
// Auto-imports are enabled, so you don't need to import components manually

// You can add reactive data here
const title = ref("{{.ProjectName}}");
</script>

<style scoped>
.container {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;
}

h1 {
  color: #42b883;
  margin-bottom: 1rem;
}

.features {
  margin: 2rem 0;
}

.features ul {
  list-style: none;
  padding: 0;
}

.features li {
  padding: 0.5rem 0;
  border-bottom: 1px solid #eee;
}

.getting-started {
  margin-top: 2rem;
  padding: 1rem;
  background: #f5f5f5;
  border-radius: 8px;
}

code {
  background: #e1e1e1;
  padding: 0.2rem 0.4rem;
  border-radius: 4px;
}

a {
  color: #42b883;
  text-decoration: none;
}

a:hover {
  text-decoration: underline;
}
</style>
