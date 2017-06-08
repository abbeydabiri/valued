until ./Valued.elf -port=83 -log=true; do
    echo "Server 'Valued.elf' crashed with exit code $?.  Respawning.." >&2
    sleep 10
done

