<script lang="ts">
	import { fade } from 'svelte/transition';
	import { goto } from '$app/navigation';
	import Sugarcube from '../components/Sugarcube.svelte';
	import { status } from '$lib/status';
	import { login } from '$lib/auth';

	let error = '';

	$: if ($status) {
		if ($status.sessionValid) {
			goto('/app');
		} else if (!$status.setupCompleted) {
			goto('/setup');
		}
	}

	function onSubmit(e: SubmitEvent) {
		const formData = new FormData(e.target as HTMLFormElement);
		const formObj = Object.fromEntries(formData.entries());
		login({
			password: formObj['password'].toString()
		})
			.then(() => {
				goto('/app');
			})
			.catch((reason) => {
				error = reason;
			});
	}
</script>

<div class="flex min-h-screen items-center justify-center">
	<div class=" border border-stone-900 dark:border-stone-100">
		<div class="p-8">
			<div class="mx-auto mb-4 w-32">
				<Sugarcube style="stroke-width:4px;" />
			</div>
			<h1 class="mb-6 text-4xl font-thin italic">Welcome back!</h1>
			<p />
			<form on:submit|preventDefault={onSubmit}>
				<label class="flex flex-col gap-4">
					Please enter your password<input
						name="password"
						id="password"
						type="password"
						autocomplete="off"
						placeholder="right here"
						class="rounded-full border border-stone-900 bg-stone-100 px-4 py-1 text-xl dark:border-stone-100 dark:bg-stone-900"
					/></label
				>
			</form>
		</div>

		{#if error}
			<p transition:fade class="bg-red-500 p-4 text-lg italic text-stone-100">
				{error}
			</p>
		{/if}
	</div>
</div>
