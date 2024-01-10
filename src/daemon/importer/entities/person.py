import xml.etree.ElementTree as ET


class Person:

    def __init__(self, contactPerson, contact):
        Person.counter += 1
        self._id = Person.counter
        self._contactPerson = contactPerson
        self._contact = contact

    def to_xml(self):
        el = ET.Element("Person")
        el.set("id", str(self._id))
        el.set("contactPerson", self._contactPerson)
        el.set("contact", self._contact)
        return el

    def get_id(self):
        return self._id

    def __str__(self):
        return f"contact:{self._contact}, contactPerson:{self._contactPerson}, id:{self._id}"


Person.counter = 0
