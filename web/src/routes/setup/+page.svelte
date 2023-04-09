<script lang="ts">
	import { fade } from 'svelte/transition';
	import { goto } from '$app/navigation';
	import { signup, status } from '$lib/stores';
	import Sugarcube from '../../components/Sugarcube.svelte';
	let error = '';

	$: if ($status) {
		if ($status.sessionValid) {
			goto('/app');
		} else if ($status.setupCompleted) {
			goto('/');
		}
	}

	function onSubmit(e: SubmitEvent) {
		const formData = new FormData(e.target as HTMLFormElement);
		const formObj = Object.fromEntries(formData.entries());

		if (!formObj['password']) {
			error = 'Please enter a password in both inputs.';
			return;
		}

		if (formObj['password'] != formObj['confirmpassword']) {
			error = 'The passwords entered do not match.';
			return;
		}

		signup(formObj['password'].toString())
			.then(() => {
				goto('/app');
				error = '';
			})
			.catch((reason) => {
				error = reason;
			});
	}
</script>

<div class="flex min-h-screen items-center justify-center">
	<div class="max-w-sm border border-stone-900 dark:border-stone-100">
		<div class="p-8">
			<div class="mx-auto mb-4 w-32">
				<Sugarcube style="stroke-width:4px;" />
			</div>
			<h1 class="mb-6 text-4xl font-thin italic">Nice to <br />meet you!</h1>
			<p />
			<form on:submit|preventDefault={onSubmit} class="flex flex-col gap-4">
				<label class="flex flex-col gap-2">
					Please choose a password<input
						name="password"
						id="password"
						type="password"
						autocomplete="off"
						placeholder="password"
						class="rounded-full border border-stone-900 bg-stone-100 px-4 py-1 text-xl dark:border-stone-100 dark:bg-stone-900"
					/></label
				>
				<label class="flex flex-col gap-2">
					And type it again<input
						name="confirmpassword"
						id="confirmpassword"
						type="password"
						autocomplete="off"
						placeholder="right here"
						class="rounded-full border border-stone-900 bg-stone-100 px-4 py-1 text-xl dark:border-stone-100 dark:bg-stone-900"
					/></label
				>

				<button
					type="submit"
					class="mt-2 rounded-full bg-stone-900 px-4 py-2 text-xl dark:bg-stone-100 dark:text-stone-900"
					>Submit</button
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
