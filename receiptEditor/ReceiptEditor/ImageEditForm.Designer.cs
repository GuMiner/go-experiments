namespace ReceiptEditor
{
    partial class ImageEditForm
    {
        /// <summary>
        /// Required designer variable.
        /// </summary>
        private System.ComponentModel.IContainer components = null;

        /// <summary>
        /// Clean up any resources being used.
        /// </summary>
        /// <param name="disposing">true if managed resources should be disposed; otherwise, false.</param>
        protected override void Dispose(bool disposing)
        {
            if (disposing && (components != null))
            {
                components.Dispose();
            }
            base.Dispose(disposing);
        }

        #region Windows Form Designer generated code

        /// <summary>
        /// Required method for Designer support - do not modify
        /// the contents of this method with the code editor.
        /// </summary>
        private void InitializeComponent()
        {
            this.brightnessTrackbar = new System.Windows.Forms.TrackBar();
            this.brightnessLabel = new System.Windows.Forms.Label();
            this.contrastLabel = new System.Windows.Forms.Label();
            this.contrastTrackBar = new System.Windows.Forms.TrackBar();
            this.categoryBox = new System.Windows.Forms.ListBox();
            this.addCategoryButton = new System.Windows.Forms.Button();
            this.newCategoryName = new System.Windows.Forms.TextBox();
            this.saveButton = new System.Windows.Forms.Button();
            this.newCategoryFolder = new System.Windows.Forms.TextBox();
            this.label1 = new System.Windows.Forms.Label();
            this.label3 = new System.Windows.Forms.Label();
            this.rotateButton = new System.Windows.Forms.Button();
            ((System.ComponentModel.ISupportInitialize)(this.brightnessTrackbar)).BeginInit();
            ((System.ComponentModel.ISupportInitialize)(this.contrastTrackBar)).BeginInit();
            this.SuspendLayout();
            // 
            // brightnessTrackbar
            // 
            this.brightnessTrackbar.LargeChange = 10;
            this.brightnessTrackbar.Location = new System.Drawing.Point(235, 12);
            this.brightnessTrackbar.Maximum = 100;
            this.brightnessTrackbar.Name = "brightnessTrackbar";
            this.brightnessTrackbar.Size = new System.Drawing.Size(164, 45);
            this.brightnessTrackbar.TabIndex = 0;
            this.brightnessTrackbar.TickFrequency = 5;
            this.brightnessTrackbar.Value = 50;
            this.brightnessTrackbar.Scroll += new System.EventHandler(this.brightnessTrackbar_Scroll);
            // 
            // brightnessLabel
            // 
            this.brightnessLabel.AutoSize = true;
            this.brightnessLabel.Location = new System.Drawing.Point(187, 15);
            this.brightnessLabel.Name = "brightnessLabel";
            this.brightnessLabel.Size = new System.Drawing.Size(56, 13);
            this.brightnessLabel.TabIndex = 1;
            this.brightnessLabel.Text = "Brightness";
            // 
            // contrastLabel
            // 
            this.contrastLabel.AutoSize = true;
            this.contrastLabel.Location = new System.Drawing.Point(187, 44);
            this.contrastLabel.Name = "contrastLabel";
            this.contrastLabel.Size = new System.Drawing.Size(46, 13);
            this.contrastLabel.TabIndex = 2;
            this.contrastLabel.Text = "Contrast";
            // 
            // contrastTrackBar
            // 
            this.contrastTrackBar.LargeChange = 10;
            this.contrastTrackBar.Location = new System.Drawing.Point(235, 44);
            this.contrastTrackBar.Maximum = 100;
            this.contrastTrackBar.Name = "contrastTrackBar";
            this.contrastTrackBar.Size = new System.Drawing.Size(164, 45);
            this.contrastTrackBar.TabIndex = 3;
            this.contrastTrackBar.TickFrequency = 5;
            this.contrastTrackBar.Value = 50;
            this.contrastTrackBar.Scroll += new System.EventHandler(this.contrastTrackBar_Scroll);
            // 
            // categoryBox
            // 
            this.categoryBox.Anchor = ((System.Windows.Forms.AnchorStyles)(((System.Windows.Forms.AnchorStyles.Top | System.Windows.Forms.AnchorStyles.Bottom) 
            | System.Windows.Forms.AnchorStyles.Right)));
            this.categoryBox.FormattingEnabled = true;
            this.categoryBox.Location = new System.Drawing.Point(12, 15);
            this.categoryBox.Name = "categoryBox";
            this.categoryBox.ScrollAlwaysVisible = true;
            this.categoryBox.Size = new System.Drawing.Size(175, 225);
            this.categoryBox.TabIndex = 8;
            // 
            // addCategoryButton
            // 
            this.addCategoryButton.Location = new System.Drawing.Point(319, 166);
            this.addCategoryButton.Name = "addCategoryButton";
            this.addCategoryButton.Size = new System.Drawing.Size(80, 23);
            this.addCategoryButton.TabIndex = 11;
            this.addCategoryButton.Text = "Add Category";
            this.addCategoryButton.UseVisualStyleBackColor = true;
            this.addCategoryButton.Click += new System.EventHandler(this.addCategoryButton_Click);
            // 
            // newCategoryName
            // 
            this.newCategoryName.Location = new System.Drawing.Point(224, 169);
            this.newCategoryName.Name = "newCategoryName";
            this.newCategoryName.Size = new System.Drawing.Size(89, 20);
            this.newCategoryName.TabIndex = 10;
            // 
            // saveButton
            // 
            this.saveButton.Location = new System.Drawing.Point(319, 217);
            this.saveButton.Name = "saveButton";
            this.saveButton.Size = new System.Drawing.Size(80, 23);
            this.saveButton.TabIndex = 13;
            this.saveButton.Text = "Save";
            this.saveButton.UseVisualStyleBackColor = true;
            this.saveButton.Click += new System.EventHandler(this.saveButton_Click);
            // 
            // newCategoryFolder
            // 
            this.newCategoryFolder.Location = new System.Drawing.Point(224, 191);
            this.newCategoryFolder.Name = "newCategoryFolder";
            this.newCategoryFolder.Size = new System.Drawing.Size(175, 20);
            this.newCategoryFolder.TabIndex = 14;
            // 
            // label1
            // 
            this.label1.AutoSize = true;
            this.label1.Location = new System.Drawing.Point(187, 194);
            this.label1.Name = "label1";
            this.label1.Size = new System.Drawing.Size(36, 13);
            this.label1.TabIndex = 15;
            this.label1.Text = "Folder";
            // 
            // label3
            // 
            this.label3.AutoSize = true;
            this.label3.Location = new System.Drawing.Point(187, 172);
            this.label3.Name = "label3";
            this.label3.Size = new System.Drawing.Size(35, 13);
            this.label3.TabIndex = 16;
            this.label3.Text = "Name";
            // 
            // rotateButton
            // 
            this.rotateButton.Location = new System.Drawing.Point(319, 76);
            this.rotateButton.Name = "rotateButton";
            this.rotateButton.Size = new System.Drawing.Size(80, 23);
            this.rotateButton.TabIndex = 17;
            this.rotateButton.Text = "Rotate";
            this.rotateButton.UseVisualStyleBackColor = true;
            this.rotateButton.Click += new System.EventHandler(this.rotateButton_Click);
            // 
            // ImageEditForm
            // 
            this.AutoScaleDimensions = new System.Drawing.SizeF(6F, 13F);
            this.AutoScaleMode = System.Windows.Forms.AutoScaleMode.Font;
            this.ClientSize = new System.Drawing.Size(411, 252);
            this.Controls.Add(this.rotateButton);
            this.Controls.Add(this.label3);
            this.Controls.Add(this.label1);
            this.Controls.Add(this.newCategoryFolder);
            this.Controls.Add(this.saveButton);
            this.Controls.Add(this.addCategoryButton);
            this.Controls.Add(this.newCategoryName);
            this.Controls.Add(this.categoryBox);
            this.Controls.Add(this.contrastTrackBar);
            this.Controls.Add(this.contrastLabel);
            this.Controls.Add(this.brightnessLabel);
            this.Controls.Add(this.brightnessTrackbar);
            this.Name = "ImageEditForm";
            this.Text = "ImageEditForm";
            ((System.ComponentModel.ISupportInitialize)(this.brightnessTrackbar)).EndInit();
            ((System.ComponentModel.ISupportInitialize)(this.contrastTrackBar)).EndInit();
            this.ResumeLayout(false);
            this.PerformLayout();

        }

        #endregion

        private System.Windows.Forms.TrackBar brightnessTrackbar;
        private System.Windows.Forms.Label brightnessLabel;
        private System.Windows.Forms.Label contrastLabel;
        private System.Windows.Forms.TrackBar contrastTrackBar;
        private System.Windows.Forms.ListBox categoryBox;
        private System.Windows.Forms.Button addCategoryButton;
        private System.Windows.Forms.TextBox newCategoryName;
        private System.Windows.Forms.Button saveButton;
        private System.Windows.Forms.TextBox newCategoryFolder;
        private System.Windows.Forms.Label label1;
        private System.Windows.Forms.Label label3;
        private System.Windows.Forms.Button rotateButton;
    }
}