import xml.etree.ElementTree as ET
from entities.job import Job

class JobPortal:

    def __init__(self, jobPortal):
        JobPortal.counter += 1
        self._id = JobPortal.counter
        self._jobPortal = jobPortal
        
        self._jobs = []

    def add_job(self, job: Job):
        self._jobs.append(job)

    def to_xml(self):
        el = ET.Element("JobPortal")
        el.set("id", str(self._id))
        el.set("jobPortal", self._jobPortal)
        

        jobs_el = ET.Element("Jobs")
        for job in self._jobs:
            jobs_el.append(job.to_xml())

        el.append(jobs_el)
        

        return el

    def get_id(self):
        return self._id

    def __str__(self):
        return f"jobPortal: {self._jobPortal}, id:{self._id}"


JobPortal.counter = 0
