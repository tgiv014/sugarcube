<script lang="ts">
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

<div class="flex min-h-screen items-center justify-center font-mono">
	<div class="max-w-sm border border-stone-900 p-8">
		<div class="mx-auto mb-4 w-32">
			<Sugarcube style="stroke-width:4px;" />
		</div>
		<h1 class="mb-6 text-4xl font-thin italic">Nice to <br />meet you!</h1>
		<p />
		<form on:submit|preventDefault={onSubmit} class="flex flex-col gap-4">
			<label class="flex flex-col gap-2">
				Please choose a new password<input
					name="password"
					id="password"
					type="password"
					autocomplete="off"
					placeholder="hunter2"
					class="rounded-full border border-stone-900 px-4 py-1 text-xl"
				/></label
			>
			<label class="flex flex-col gap-2">
				And type it again<input
					name="confirmpassword"
					id="confirmpassword"
					type="password"
					autocomplete="off"
					placeholder="right here"
					class="rounded-full border border-stone-900 px-4 py-1 text-xl"
				/></label
			>

			<button type="submit" class="mt-2 rounded-full bg-stone-900 px-4 py-2 text-xl text-stone-100"
				>Submit</button
			>
			{#if error}
				<p class="bg-red-700 p-4 font-bold text-stone-100">{error}</p>
			{/if}
		</form>
	</div>
</div>
