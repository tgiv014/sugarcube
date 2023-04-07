<script lang="ts">
	import { goto } from '$app/navigation';
	import { signup, status } from '$lib/stores';
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
	<div class="border border-stone-900 p-8 max-w-sm">
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

			<button type="submit" class="rounded-full bg-stone-900 px-4 py-2 mt-2 text-stone-100 text-xl"
				>Submit</button
			>
			{#if error}
				<p class="bg-red-700 text-stone-100 font-bold p-4">{error}</p>
			{/if}
		</form>
	</div>
</div>
