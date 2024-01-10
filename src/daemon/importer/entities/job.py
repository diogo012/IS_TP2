import xml.etree.ElementTree as ET


class Job:

    def __init__(self, jobDescription, jobTitle, experience, workType, qualifications, skills, preference, jobPostingDate, personId, role):
        Job.counter += 1
        self._id = Job.counter
        self._jobDescription = jobDescription
        self._jobTitle = jobTitle
        self._experience = experience
        self._workType = workType
        self._qualifications = qualifications
        self._skills = skills
        self._preference = preference
        self._jobPostingDate = jobPostingDate
        self._personId = personId
        self._role = role

    def to_xml(self):
        el = ET.Element("Job")
        el.set("id", str(self._id))
        el.set("jobTitle", self._jobTitle)
        el.set("experience", self._experience)
        el.set("workType", self._workType)
        el.set("qualifications", self._qualifications)
        el.set("preference", self._preference)
        el.set("jobPostingDate", self._jobPostingDate)
        el.set("person_ref", str(self._personId.get_id()))
        el.set("role_ref", str(self._role.get_id()))

        
        # Adicionando cada jobDescription como um elemento
        desc_el = ET.SubElement(el, "Description")
        desc_el.text = self._jobDescription

        # Adicionando cada skills como um elemento 
        skills_el = ET.SubElement(el, "Skills")
        skills_el.text = self._skills


        return el

    def __str__(self):
        return f"jobTitle:{self._jobTitle}, person_ref: {self._personId}, jobPostingDate:{self._jobPostingDate}, preference: {self._preference}, skills:{self._skills}, qualifications: {self._qualifications}, workType: {self._workType}, experience: {self._experience}, jobDescription: {self._jobDescription}, role_ref:{self._role}, id:{self._id}"


Job.counter = 0
