<script>
	import '../app.css';
	import { goto } from '$app/navigation';
	import { status } from '../lib/stores';
	import { page } from '$app/stores';
	$: console.log($status);
	$: {
		if ($status && !$status.setupCompleted && $page.url.pathname !== '/setup') {
			goto('/setup');
		}
	}
</script>

<slot />
