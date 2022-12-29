<template>
  <div class="item">
    <h3>
      <slot name="name"></slot>
    </h3>
    <div class="editor">
      <h3>
        <slot name="editor"></slot>
      </h3>
      <slot></slot>
    </div>
  </div>
</template>
