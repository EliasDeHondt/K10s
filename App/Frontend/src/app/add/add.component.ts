/**********************************/
/* @since 01/01/2025              */
/* @author K10s Open Source Team  */
/**********************************/

import { Component } from '@angular/core';
import { NavComponent } from "../nav/nav.component";
import { FooterComponent } from "../footer/footer.component";
import { FormsModule } from "@angular/forms";
import { DeploymentService } from "../services/deployment.service";
import { TranslatePipe } from "@ngx-translate/core";

@Component({
    selector: 'app-add-deployment',
    templateUrl: './add.component.html',
    styleUrl: './add.component.css',
    imports: [NavComponent, FooterComponent, FormsModule, TranslatePipe],
    standalone: true
})

export class AddComponent {
    yamlText: string = '';
    fileContent: string | null = null;
    fileUploaded: boolean = false;
    textAreaActive: boolean = false;

    constructor(private deploymentService: DeploymentService) {}

    onFileUpload(event: Event) {
        const input = event.target as HTMLInputElement;
        if (input.files && input.files.length > 0) {
            const file = input.files[0];

            const reader = new FileReader();
            reader.onload = (e) => {
                this.fileContent = e.target?.result as string;
                this.fileUploaded = true;
                this.clearTextarea();
                this.textAreaActive = false;
            };
            reader.readAsText(file);
        }
    }

    onTextInput() {
        this.fileUploaded = false;
        this.fileContent = null;
        this.textAreaActive = true;
    }

    sendData() {
        const yamlData = this.fileContent || this.yamlText;

        if (!yamlData) {
            console.warn('No YAML data to upload');
            return;
        }

        this.deploymentService.uploadYaml(yamlData).subscribe({
            next: (response) => {console.log('YAML upload success', response);  this.clearTextarea();},
            error: (error) => console.error('Upload failed', error),
        });
    }

    clearTextarea() {
        this.yamlText = '';
        this.fileUploaded = true;
        this.fileContent = null;
        this.textAreaActive = false;
    }
}