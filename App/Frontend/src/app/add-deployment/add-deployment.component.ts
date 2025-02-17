import { Component } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {NavComponent} from "../nav/nav.component";
import {FooterComponent} from "../footer/footer.component";
import {FormsModule} from "@angular/forms";

@Component({
  selector: 'app-add-deployment',
  imports: [
    NavComponent,
    FooterComponent,
    FormsModule
  ],
  templateUrl: './add-deployment.component.html',
  standalone: true,
  styleUrl: './add-deployment.component.css'
})
export class AddDeploymentComponent {
  yamlText: string = '';
  fileName: string = '';
  isFileSelected = false;

  constructor(private http: HttpClient) {}

  onFileSelected(event: any) {
    const file = event.target.files[0];
    if (file) {
      this.fileName = file.name;
      this.isFileSelected = true;
      this.yamlText = ''; // Clear text area
      const reader = new FileReader();
      reader.onload = () => {
        this.yamlText = reader.result as string;
      };
      reader.readAsText(file);
    }
  }

  onTextChange() {
    if (this.yamlText.trim() !== '') {
      this.isFileSelected = false;
      this.fileName = '';
    }
  }

  submitYaml() {
    if (!this.yamlText.trim()) {
      alert('Please provide YAML data.');
      return;
    }

    this.http.post('/api/yaml/upload', { yaml: this.yamlText }).subscribe(response => {
      console.log('YAML uploaded successfully', response);
    });
  }
}
