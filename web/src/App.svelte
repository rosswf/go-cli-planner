<script>
  import { onMount } from "svelte";
  import Task from "./Task.svelte";

  let tasks = [];
  let newTask = "";

  onMount(async () => {
    const res = await fetch("http://localhost:5000/tasks");
    tasks = await res.json();
    console.log(tasks);
  });

  async function addTask() {
    const res = await fetch("http://localhost:5000/tasks", {
      method: "POST",
      body: `{ "Name": "${newTask}" }`,
    });
    const addedTask = await res.json();
    console.log(newTask);
    tasks = [...tasks, addedTask[0]];
    newTask = "";
  }
</script>

<article class="to-do-list">
  <header>My Awesome To-Do List</header>
  <table role="grid">
    {#each tasks as task (task)}
      <Task {task} />
    {/each}
  </table>
  <input type="text" bind:value={newTask} />
  <button on:click={addTask}>Add</button>
</article>

<style>
</style>
