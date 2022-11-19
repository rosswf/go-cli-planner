<script>
    export let task;
    import { fly } from "svelte/transition";

    async function completeTask(event) {
        await fetch("http://localhost:5000/tasks/" + event.target.id, {
            method: "POST",
        });
    }
</script>

<tr
    in:fly={{ delay: 0, duration: 700, x: -200, y: 0 }}
    out:fly={{ delay: 0, duration: 700, x: -200, y: 0 }}
>
    <td
        ><input
            type="checkbox"
            id={task.id}
            bind:checked={task.complete}
            on:click={completeTask}
        />{task.name}</td
    >
</tr>
