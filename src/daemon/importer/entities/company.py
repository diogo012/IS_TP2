import xml.etree.ElementTree as ET

class Company:

    def __init__(self, company, companySize, benefits, country):
        Company.counter += 1
        self._id = Company.counter
        self._company = company
        self._companySize = companySize
        self._benefits = benefits
        self._country = country

    def to_xml(self):
        el = ET.Element("Company")
        el.set("id", str(self._id))
        el.set("company", self._company)
        el.set("companySize", self._companySize)
        el.set("country_ref", str(self._country.get_id()))

        cleaned_benefits = self._benefits.replace('{', '').replace('}', '').replace("'",'')

        benefits_el = ET.SubElement(el, "Benefits")
        benefits_el.text = cleaned_benefits
        return el
    
    def get_id(self):
        return self._id

    def __str__(self):
        return f"country:{self._country}, benefits:{self._benefits}, companySize:{self._companySize}, company:{self._company}, id:{self._id}"


Company.counter = 0
