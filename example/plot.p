set title "Bandwith pod vs. vm"
set terminal pngcairo size 1250,1062 enhanced font 'Verdana,10'
set xlabel "Time"
set ylabel "Bandwith"
set output 'pod-vs-vm-bandwith-randwrite-4k.png'
plot 'vm/fio-4k-device-to-test-write-seq.results_bw.1.log' with lines title 'vm-4k', \
        'pod/fio-4k-device-to-test-write-seq.results_bw.1.log' with lines title 'pod-4k'

set output 'pod-vs-vm-bandwith-randwrite-128k.png'
plot 'vm/fio-128k-device-to-test-write-seq.results_bw.2.log' with lines title 'vm-128k', \
        'pod/fio-128k-device-to-test-write-seq.results_bw.2.log' with lines title 'pod-128k'

