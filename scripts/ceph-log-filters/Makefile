HOSTS= 10.92.100.1 10.92.100.2 10.92.100.3 10.92.100.4 \
       10.121.100.1 10.121.100.2 10.121.100.3 10.121.100.4

# this target is for daily cron-jobs
cron: 
	make refresh >/dev/null 2>&1
	make filter
refresh: refresh_radosgw refresh_audit  refresh_ceph  refresh_mgr refresh_mon refresh_osd refresh_admin 
filter: radosgw audit ceph mgr mon osd admin

# This target is for debugging
all: admin



#
# radosgw -- the non-empty radosgw files
#
radosgw:  ./radosgw.filter 
	@awk -f ./radosgw.filter ./radosgw.log

radosgw.filter: radosgw.template
	mkfilter radosgw.template

#
# audit
#
audit:  ./audit.filter
	@awk -f ./audit.filter ./ceph.audit.log

audit.filter: audit.template
	mkfilter audit.template

#
# ceph
#
ceph:  ./ceph.filter
	@awk -f ./ceph.filter ./ceph.log

ceph.filter: ceph.template
	mkfilter ceph.template

#
# mgr
#
mgr:  ./mgr.filter
	@awk -f ./mgr.filter ./mgr.log

mgr.filter: mgr.template
	mkfilter mgr.template

#
# mon
#
mon:  ./mon.filter
	@awk -f ./mon.filter ./mon.log

mon.filter: mon.template
	mkfilter mon.template

#
# osd
#
osd:  ./osd.filter
	@awk -f ./osd.filter ./osd.log

osd.filter: osd.template
	mkfilter osd.template

#
# admin
#
admin:  ./admin.filter
	@awk -f ./admin.filter ./admin.log

admin.filter: admin.template
	mkfilter admin.template


# The extra call to ssh to mac is because of a routing bug
# other prople shouldn't need to do so.

refresh_radosgw:
	-for i in ${HOSTS}; do \
		ssh -2 -A -Y dcollierbrown@mac \
		ssh -2 -A dcollierbrown@$${i} \
		sudo cat /var/log/ceph/ceph.client.radosgw.log; \
	done >radosgw.log

refresh_audit:
	-for i in ${HOSTS}; do \
		ssh -2 -A -Y dcollierbrown@mac \
		ssh -2 -A dcollierbrown@$${i} \
		sudo cat /var/log/ceph/ceph.audit.log; \
	done >ceph.audit.log

refresh_ceph:
	-for i in ${HOSTS}; do \
		ssh -2 -A -Y dcollierbrown@mac \
		ssh -2 -A dcollierbrown@$${i} \
		sudo cat /var/log/ceph/ceph.log; \
	done >ceph.log

refresh_mgr:
	-for i in ${HOSTS}; do \
		for j in t c; do \
			for k in 1 2 3 4; do \
				ssh -2 -A -Y dcollierbrown@mac \
				ssh -2 -A dcollierbrown@$${i} \
				sudo cat /var/log/ceph/ceph-mgr.p$${j}-ceph-00$${k}.log; \
			done; \
		done; \
	done > mgr.log

refresh_mon:
	-for i in ${HOSTS}; do \
		for j in t c; do \
			for k in 1 2 3 4; do \
				ssh -2 -A -Y dcollierbrown@mac \
				ssh -2 -A dcollierbrown@$${i} \
				sudo cat /var/log/ceph/ceph-mon.p$${j}-ceph-00$${k}.log; \
			done; \
		done; \
	done > mon.log

OSDS=0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 19 20 21 22 23 24
refresh_osd:
	-for i in ${HOSTS}; do \
		for k in ${OSDS}; do \
			ssh -2 -A -Y dcollierbrown@mac \
			ssh -2 -A dcollierbrown@$${i} \
			sudo cat /var/log/ceph/ceph-osd.$${k}.log; \
		done; \
	done > osd.log

refresh_admin:
	-for i in ${HOSTS}; do \
		ssh -2 -A -Y dcollierbrown@mac \
		ssh -2 -A dcollierbrown@$${i} \
		sudo cat /var/log/ceph/ceph-osd.admin.log; \
	done > admin.log

